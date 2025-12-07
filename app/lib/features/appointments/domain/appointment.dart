import 'package:freezed_annotation/freezed_annotation.dart';

// ignore_for_file: invalid_annotation_target

part 'appointment.freezed.dart';
part 'appointment.g.dart';

@freezed
abstract class Appointment with _$Appointment {
  const factory Appointment({
    @JsonKey(name: '_id') required String id,
    required DateTime startTime,
    required DateTime endTime,
    required String status,
    Map<String, dynamic>? pt,
    Map<String, dynamic>? patient,
  }) = _Appointment;

  factory Appointment.fromJson(Map<String, dynamic> json) =>
      _$AppointmentFromJson(json);
}
