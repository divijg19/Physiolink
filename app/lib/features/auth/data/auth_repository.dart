import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/api/api_client.dart';

part 'auth_repository.g.dart';

@Riverpod(keepAlive: true)
AuthRepository authRepository(Ref ref) {
  return AuthRepository(ref.watch(apiClientProvider));
}

class AuthRepository {
  final Dio _dio;

  AuthRepository(this._dio);

  Future<String> login(String email, String password) async {
    final response = await _dio.post(
      '/auth/login',
      data: {'email': email, 'password': password},
    );
    return response.data['token'];
  }

  Future<void> register(String email, String password, String role) async {
    await _dio.post(
      '/auth/register',
      data: {'email': email, 'password': password, 'role': role},
    );
  }

  /// Returns the signed-in user's email from `GET /profile/me`.
  ///
  /// The JWT intentionally carries only `user.id` and `user.role`, so that no
  /// PII is baked into a token that may be logged or cached. The email has to
  /// be read from the profile endpoint instead.
  Future<String> fetchCurrentUserEmail() async {
    final response = await _dio.get('/profile/me');
    final data = response.data;
    if (data is Map && data['user'] is Map) {
      final email = (data['user'] as Map)['email'];
      if (email is String) return email;
    }
    return '';
  }
}
