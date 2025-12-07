// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'therapist.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_Therapist _$TherapistFromJson(Map<String, dynamic> json) => _Therapist(
  id: json['_id'] as String,
  email: json['email'] as String,
  availableSlotsCount: (json['availableSlotsCount'] as num?)?.toInt() ?? 0,
  reviewCount: (json['reviewCount'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$TherapistToJson(_Therapist instance) =>
    <String, dynamic>{
      '_id': instance.id,
      'email': instance.email,
      'availableSlotsCount': instance.availableSlotsCount,
      'reviewCount': instance.reviewCount,
    };
