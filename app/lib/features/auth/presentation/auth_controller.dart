import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/api/api_client.dart';
import '../data/auth_repository.dart';
import '../domain/user.dart';

part 'auth_controller.g.dart';

@Riverpod(keepAlive: true)
class AuthController extends _$AuthController {
  @override
  Future<User?> build() async {
    final storage = ref.watch(secureStorageProvider);
    final token = await storage.read(key: 'auth_token');
    if (token == null) return null;
    return _restoreSession(token, storage);
  }

  Future<void> login(String email, String password) async {
    state = const AsyncValue.loading();
    try {
      final repo = ref.read(authRepositoryProvider);
      final token = await repo.login(email, password);

      final storage = ref.read(secureStorageProvider);
      await storage.write(key: 'auth_token', value: token);

      state = AsyncValue.data(await _resolveUser(token));
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> logout() async {
    final storage = ref.read(secureStorageProvider);
    await storage.delete(key: 'auth_token');
    state = const AsyncValue.data(null);
  }

  Future<void> register(String email, String password, String role) async {
    state = const AsyncValue.loading();
    try {
      final repo = ref.read(authRepositoryProvider);
      await repo.register(email, password, role);
      // Auto login after register? Or just redirect to login?
      // For now, let's just stop loading and let UI handle redirect
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  /// Builds a [User] from an already-validated [token].
  User _userFromToken(String token) {
    return userFromTokenClaims(JwtDecoder.decode(token));
  }

  /// Decodes [token] and enriches it with the user's email.
  Future<User> _resolveUser(String token) async {
    final user = _userFromToken(token);
    try {
      final email = await ref
          .read(authRepositoryProvider)
          .fetchCurrentUserEmail();
      if (email.isNotEmpty) return user.copyWith(email: email);
    } catch (_) {
      // The email is cosmetic; a profile hiccup must not break the session.
    }
    return user;
  }

  /// Restores a session from a stored [token], discarding the token if it is
  /// expired or unparseable so a corrupt value cannot wedge the app on a
  /// loading spinner forever.
  Future<User?> _restoreSession(
    String token,
    FlutterSecureStorage storage,
  ) async {
    try {
      if (JwtDecoder.isExpired(token)) {
        await storage.delete(key: 'auth_token');
        return null;
      }
      return await _resolveUser(token);
    } catch (_) {
      await storage.delete(key: 'auth_token');
      return null;
    }
  }
}
