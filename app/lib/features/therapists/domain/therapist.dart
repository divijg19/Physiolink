import 'package:freezed_annotation/freezed_annotation.dart';

// ignore_for_file: invalid_annotation_target

part 'therapist.freezed.dart';
part 'therapist.g.dart';

@freezed
abstract class Therapist with _$Therapist {
  const factory Therapist({
    @JsonKey(name: '_id') required String id,
    required String email,
    @Default(0) int availableSlotsCount,
    @Default(0) int reviewCount,
    // Add profile fields if needed
  }) = _Therapist;

  factory Therapist.fromJson(Map<String, dynamic> json) =>
      _$TherapistFromJson(json);
}
