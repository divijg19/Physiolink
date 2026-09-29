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

  /// Books [slotId] and refreshes the caller's schedule.
  ///
  /// Errors are recorded in [state] *and* rethrown so the calling widget can
  /// report the failure. Swallowing them here made the UI announce
  /// "Appointment booked!" for bookings that never happened.
  Future<void> bookAppointment(String slotId) async {
    state = const AsyncLoading();
    try {
      final repo = ref.read(appointmentRepositoryProvider);
      await repo.bookAppointment(slotId);
      ref.invalidate(myAppointmentsProvider);
      state = const AsyncData(null);
    } catch (error, stackTrace) {
      state = AsyncError(error, stackTrace);
      rethrow;
    }
  }
}
