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

    if (token != null && !JwtDecoder.isExpired(token)) {
      final decoded = JwtDecoder.decode(token);
      // Assuming the token contains user info, otherwise fetch profile
      return User(
        id: decoded['sub'] ?? '',
        email: decoded['email'] ?? '',
        role: decoded['role'] ?? 'patient',
      );
    }
    return null;
  }

  Future<void> login(String email, String password) async {
    state = const AsyncValue.loading();
    try {
      final repo = ref.read(authRepositoryProvider);
      final token = await repo.login(email, password);

      final storage = ref.read(secureStorageProvider);
      await storage.write(key: 'auth_token', value: token);

      final decoded = JwtDecoder.decode(token);
      final user = User(
        id: decoded['sub'] ?? '',
        email: decoded['email'] ?? '',
        role: decoded['role'] ?? 'patient',
      );

      state = AsyncValue.data(user);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
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
}
