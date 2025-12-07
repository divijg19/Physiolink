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

  Future<Appointment> bookAppointment(String slotId) async {
    final response = await _dio.put('/appointments/$slotId/book');
    return Appointment.fromJson(response.data);
  }
}
