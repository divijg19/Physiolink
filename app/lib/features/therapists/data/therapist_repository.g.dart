// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'therapist_repository.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(therapistRepository)
const therapistRepositoryProvider = TherapistRepositoryProvider._();

final class TherapistRepositoryProvider
    extends
        $FunctionalProvider<
          TherapistRepository,
          TherapistRepository,
          TherapistRepository
        >
    with $Provider<TherapistRepository> {
  const TherapistRepositoryProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'therapistRepositoryProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$therapistRepositoryHash();

  @$internal
  @override
  $ProviderElement<TherapistRepository> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  TherapistRepository create(Ref ref) {
    return therapistRepository(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(TherapistRepository value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<TherapistRepository>(value),
    );
  }
}

String _$therapistRepositoryHash() =>
    r'7de4011c9e7782a6059117749ae6083f6d984561';
