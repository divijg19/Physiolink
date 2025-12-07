// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'therapist_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(therapists)
const therapistsProvider = TherapistsProvider._();

final class TherapistsProvider
    extends
        $FunctionalProvider<
          AsyncValue<List<Therapist>>,
          List<Therapist>,
          FutureOr<List<Therapist>>
        >
    with $FutureModifier<List<Therapist>>, $FutureProvider<List<Therapist>> {
  const TherapistsProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'therapistsProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$therapistsHash();

  @$internal
  @override
  $FutureProviderElement<List<Therapist>> $createElement(
    $ProviderPointer pointer,
  ) => $FutureProviderElement(pointer);

  @override
  FutureOr<List<Therapist>> create(Ref ref) {
    return therapists(ref);
  }
}

String _$therapistsHash() => r'c75a13b7e9c6d092c411c6e2a90db3f77afca647';
