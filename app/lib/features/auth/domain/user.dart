import 'package:freezed_annotation/freezed_annotation.dart';

// ignore_for_file: invalid_annotation_target

part 'user.freezed.dart';
part 'user.g.dart';

@freezed
abstract class User with _$User {
  const factory User({
    @JsonKey(name: '_id') required String id,
    required String email,
    required String role, // 'patient' or 'therapist'
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

/// Builds a [User] from the decoded claims of a JWT issued by the Go backend.
///
/// The backend signs `{"user": {"id": ..., "role": ...}, "exp": ...}` — the
/// claims are *nested*, and there is deliberately no flat `sub`/`email`/`role`
/// and no email in the token at all (no PII in a token that may be logged or
/// cached). Reading flat keys yields an empty id and a hardcoded "patient" role
/// for every account, so the shape is asserted in user_token_test.dart.
User userFromTokenClaims(Map<String, dynamic> claims) {
  final nested = claims['user'];
  final user = nested is Map ? nested : const <String, dynamic>{};
  return User(
    id: (user['id'] as String?) ?? '',
    // Populated separately from GET /profile/me; the token never carries it.
    email: '',
    role: (user['role'] as String?) ?? 'patient',
  );
}
