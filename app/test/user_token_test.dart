import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:physiolink_app/features/auth/domain/user.dart';

/// These fixtures are byte-for-byte the payload shape produced by the Go
/// backend in `internal/handlers/auth.go`:
///
///     jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
///         "user": map[string]interface{}{"id": id, "role": role},
///         "exp": time.Now().Add(5 * time.Hour).Unix(),
///     })
///
/// If the backend claim shape changes, update these and the Dart parser
/// together — decoding flat `sub`/`email`/`role` keys silently produced an
/// empty id and a "patient" role for every account.
Map<String, dynamic> _claimsWith(Map<String, dynamic> user) => {
      'exp': 1790710499,
      'user': user,
    };

void main() {
  group('userFromTokenClaims', () {
    test('reads the nested user id and role', () {
      final user = userFromTokenClaims(
        _claimsWith({'id': '550e8400-e29b-41d4-a716-446655440000', 'role': 'pt'}),
      );

      expect(user.id, '550e8400-e29b-41d4-a716-446655440000');
      expect(user.role, 'pt');
    });

    test('preserves a patient role', () {
      final user = userFromTokenClaims(
        _claimsWith({'id': 'u-2', 'role': 'patient'}),
      );

      expect(user.role, 'patient');
    });

    test('never fabricates an email (it is not in the token)', () {
      final user = userFromTokenClaims(_claimsWith({'id': 'u-3', 'role': 'pt'}));

      expect(user.email, isEmpty);
    });

    test('falls back to the patient role when the claim is absent', () {
      // Documented fallback: a token without a role claim yields "patient".
      // The id must still be read correctly from the nested claim.
      final user = userFromTokenClaims({'user': {'id': 'u-4'}});

      expect(user.id, 'u-4');
      expect(user.role, 'patient');
    });

    test('tolerates a token with no user claim', () {
      final user = userFromTokenClaims({'exp': 1790710499});

      expect(user.id, isEmpty);
      expect(user.role, 'patient');
    });

    test('ignores flat sub/role keys, which the backend never sets', () {
      // Guards the original defect: reading these returned the flat values and
      // masked the nested ones, so role always looked like "patient".
      final user = userFromTokenClaims({
        'sub': 'flat-id',
        'role': 'admin',
        'user': {'id': 'nested-id', 'role': 'pt'},
      });

      expect(user.id, 'nested-id');
      expect(user.role, 'pt');
    });

    test('matches the real base64 payload emitted by the backend', () {
      // {"exp":1790710499,"user":{"id":"u-1","role":"pt"}}
      const payload =
          'eyJleHAiOjE3OTA3MTA0OTksInVzZXIiOnsiaWQiOiJ1LTEiLCJyb2xlIjoicHQifX0';
      final decoded = jsonDecode(
        utf8.decode(base64Url.decode(base64Url.normalize(payload))),
      ) as Map<String, dynamic>;

      final user = userFromTokenClaims(decoded);

      expect(user.id, 'u-1');
      expect(user.role, 'pt');
    });
  });
}
