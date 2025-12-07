import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../data/appointment_repository.dart';
import '../domain/appointment.dart';

part 'appointment_controller.g.dart';

@riverpod
Future<List<Appointment>> myAppointments(Ref ref) async {
  final repo = ref.watch(appointmentRepositoryProvider);
  return repo.getMyAppointments();
}

@riverpod
Future<List<Appointment>> therapistAvailability(Ref ref, String ptId) async {
  final repo = ref.watch(appointmentRepositoryProvider);
  return repo.getTherapistAvailability(ptId);
}

@riverpod
class AppointmentController extends _$AppointmentController {
  @override
  FutureOr<void> build() {
    // no-op
  }

  Future<void> bookAppointment(String slotId) async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      final repo = ref.read(appointmentRepositoryProvider);
      await repo.bookAppointment(slotId);
      // Invalidate my appointments to refresh the list
      ref.invalidate(myAppointmentsProvider);
    });
  }
}
