// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'appointment_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(myAppointments)
final myAppointmentsProvider = MyAppointmentsProvider._();

final class MyAppointmentsProvider
    extends
        $FunctionalProvider<
          AsyncValue<List<Appointment>>,
          List<Appointment>,
          FutureOr<List<Appointment>>
        >
    with
        $FutureModifier<List<Appointment>>,
        $FutureProvider<List<Appointment>> {
  MyAppointmentsProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'myAppointmentsProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$myAppointmentsHash();

  @$internal
  @override
  $FutureProviderElement<List<Appointment>> $createElement(
    $ProviderPointer pointer,
  ) => $FutureProviderElement(pointer);

  @override
  FutureOr<List<Appointment>> create(Ref ref) {
    return myAppointments(ref);
  }
}

String _$myAppointmentsHash() => r'0bb408f0e3e4c636626c0ccca72d1c4d7bb72af1';

@ProviderFor(therapistAvailability)
final therapistAvailabilityProvider = TherapistAvailabilityFamily._();

final class TherapistAvailabilityProvider
    extends
        $FunctionalProvider<
          AsyncValue<List<Appointment>>,
          List<Appointment>,
          FutureOr<List<Appointment>>
        >
    with
        $FutureModifier<List<Appointment>>,
        $FutureProvider<List<Appointment>> {
  TherapistAvailabilityProvider._({
    required TherapistAvailabilityFamily super.from,
    required String super.argument,
  }) : super(
         retry: null,
         name: r'therapistAvailabilityProvider',
         isAutoDispose: true,
         dependencies: null,
         $allTransitiveDependencies: null,
       );

  @override
  String debugGetCreateSourceHash() => _$therapistAvailabilityHash();

  @override
  String toString() {
    return r'therapistAvailabilityProvider'
        ''
        '($argument)';
  }

  @$internal
  @override
  $FutureProviderElement<List<Appointment>> $createElement(
    $ProviderPointer pointer,
  ) => $FutureProviderElement(pointer);

  @override
  FutureOr<List<Appointment>> create(Ref ref) {
    final argument = this.argument as String;
    return therapistAvailability(ref, argument);
  }

  @override
  bool operator ==(Object other) {
    return other is TherapistAvailabilityProvider && other.argument == argument;
  }

  @override
  int get hashCode {
    return argument.hashCode;
  }
}

String _$therapistAvailabilityHash() =>
    r'8180057d0a3c98ebd2eaf7aa6a3f87fe5f122ca0';

final class TherapistAvailabilityFamily extends $Family
    with $FunctionalFamilyOverride<FutureOr<List<Appointment>>, String> {
  TherapistAvailabilityFamily._()
    : super(
        retry: null,
        name: r'therapistAvailabilityProvider',
        dependencies: null,
        $allTransitiveDependencies: null,
        isAutoDispose: true,
      );

  TherapistAvailabilityProvider call(String ptId) =>
      TherapistAvailabilityProvider._(argument: ptId, from: this);

  @override
  String toString() => r'therapistAvailabilityProvider';
}

@ProviderFor(AppointmentController)
final appointmentControllerProvider = AppointmentControllerProvider._();

final class AppointmentControllerProvider
    extends $AsyncNotifierProvider<AppointmentController, void> {
  AppointmentControllerProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'appointmentControllerProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$appointmentControllerHash();

  @$internal
  @override
  AppointmentController create() => AppointmentController();
}

String _$appointmentControllerHash() =>
    r'e84783f79ea288e9a0265795484ee7dc2c33f2e2';

abstract class _$AppointmentController extends $AsyncNotifier<void> {
  FutureOr<void> build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<AsyncValue<void>, void>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<void>, void>,
              AsyncValue<void>,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
