import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../../core/api/api_client.dart';
import '../domain/appointment.dart';

part 'appointment_repository.g.dart';

@Riverpod(keepAlive: true)
AppointmentRepository appointmentRepository(Ref ref) {
  return AppointmentRepository(ref.watch(apiClientProvider));
}

class AppointmentRepository {
  final Dio _dio;

  AppointmentRepository(this._dio);

  Future<List<Appointment>> getMyAppointments() async {
    final response = await _dio.get('/appointments/me');
    final data = response.data as List;
    return data.map((e) => Appointment.fromJson(e)).toList();
  }

  Future<List<Appointment>> getTherapistAvailability(String ptId) async {
    final response = await _dio.get('/appointments/availability/$ptId');
    final data = response.data as List;
    return data.map((e) => Appointment.fromJson(e)).toList();
  }

  /// Books [slotId] and returns the id of the created appointment.
  ///
  /// The endpoint responds with `{"id": "<uuid>"}`, not a full [Appointment],
  /// so this must not try to deserialise an [Appointment] from the body.
  Future<String> bookAppointment(String slotId) async {
    final response = await _dio.put('/appointments/$slotId/book');
    final data = response.data;
    if (data is! Map || (data['id'] == null && data['_id'] == null)) {
      throw StateError('Unexpected booking response: $data');
    }
    return (data['id'] ?? data['_id']) as String;
  }
}
