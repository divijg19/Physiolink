import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/api/api_client.dart';
import '../domain/therapist.dart';

part 'therapist_repository.g.dart';

@Riverpod(keepAlive: true)
TherapistRepository therapistRepository(Ref ref) {
  return TherapistRepository(ref.watch(apiClientProvider));
}

class TherapistRepository {
  final Dio _dio;

  TherapistRepository(this._dio);

  Future<List<Therapist>> getTherapists() async {
    final response = await _dio.get('/therapists');
    final data = response.data['data'] as List;
    return data.map((e) => Therapist.fromJson(e)).toList();
  }
}
