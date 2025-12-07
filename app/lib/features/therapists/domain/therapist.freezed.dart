// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'therapist.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$Therapist {

@JsonKey(name: '_id') String get id; String get email; int get availableSlotsCount; int get reviewCount;
/// Create a copy of Therapist
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TherapistCopyWith<Therapist> get copyWith => _$TherapistCopyWithImpl<Therapist>(this as Therapist, _$identity);

  /// Serializes this Therapist to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Therapist&&(identical(other.id, id) || other.id == id)&&(identical(other.email, email) || other.email == email)&&(identical(other.availableSlotsCount, availableSlotsCount) || other.availableSlotsCount == availableSlotsCount)&&(identical(other.reviewCount, reviewCount) || other.reviewCount == reviewCount));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,email,availableSlotsCount,reviewCount);

@override
String toString() {
  return 'Therapist(id: $id, email: $email, availableSlotsCount: $availableSlotsCount, reviewCount: $reviewCount)';
}


}

/// @nodoc
abstract mixin class $TherapistCopyWith<$Res>  {
  factory $TherapistCopyWith(Therapist value, $Res Function(Therapist) _then) = _$TherapistCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: '_id') String id, String email, int availableSlotsCount, int reviewCount
});




}
/// @nodoc
class _$TherapistCopyWithImpl<$Res>
    implements $TherapistCopyWith<$Res> {
  _$TherapistCopyWithImpl(this._self, this._then);

  final Therapist _self;
  final $Res Function(Therapist) _then;

/// Create a copy of Therapist
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? email = null,Object? availableSlotsCount = null,Object? reviewCount = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,email: null == email ? _self.email : email // ignore: cast_nullable_to_non_nullable
as String,availableSlotsCount: null == availableSlotsCount ? _self.availableSlotsCount : availableSlotsCount // ignore: cast_nullable_to_non_nullable
as int,reviewCount: null == reviewCount ? _self.reviewCount : reviewCount // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [Therapist].
extension TherapistPatterns on Therapist {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Therapist value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Therapist() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Therapist value)  $default,){
final _that = this;
switch (_that) {
case _Therapist():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Therapist value)?  $default,){
final _that = this;
switch (_that) {
case _Therapist() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: '_id')  String id,  String email,  int availableSlotsCount,  int reviewCount)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Therapist() when $default != null:
return $default(_that.id,_that.email,_that.availableSlotsCount,_that.reviewCount);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: '_id')  String id,  String email,  int availableSlotsCount,  int reviewCount)  $default,) {final _that = this;
switch (_that) {
case _Therapist():
return $default(_that.id,_that.email,_that.availableSlotsCount,_that.reviewCount);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: '_id')  String id,  String email,  int availableSlotsCount,  int reviewCount)?  $default,) {final _that = this;
switch (_that) {
case _Therapist() when $default != null:
return $default(_that.id,_that.email,_that.availableSlotsCount,_that.reviewCount);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Therapist implements Therapist {
  const _Therapist({@JsonKey(name: '_id') required this.id, required this.email, this.availableSlotsCount = 0, this.reviewCount = 0});
  factory _Therapist.fromJson(Map<String, dynamic> json) => _$TherapistFromJson(json);

@override@JsonKey(name: '_id') final  String id;
@override final  String email;
@override@JsonKey() final  int availableSlotsCount;
@override@JsonKey() final  int reviewCount;

/// Create a copy of Therapist
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TherapistCopyWith<_Therapist> get copyWith => __$TherapistCopyWithImpl<_Therapist>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TherapistToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Therapist&&(identical(other.id, id) || other.id == id)&&(identical(other.email, email) || other.email == email)&&(identical(other.availableSlotsCount, availableSlotsCount) || other.availableSlotsCount == availableSlotsCount)&&(identical(other.reviewCount, reviewCount) || other.reviewCount == reviewCount));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,email,availableSlotsCount,reviewCount);

@override
String toString() {
  return 'Therapist(id: $id, email: $email, availableSlotsCount: $availableSlotsCount, reviewCount: $reviewCount)';
}


}

/// @nodoc
abstract mixin class _$TherapistCopyWith<$Res> implements $TherapistCopyWith<$Res> {
  factory _$TherapistCopyWith(_Therapist value, $Res Function(_Therapist) _then) = __$TherapistCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: '_id') String id, String email, int availableSlotsCount, int reviewCount
});




}
/// @nodoc
class __$TherapistCopyWithImpl<$Res>
    implements _$TherapistCopyWith<$Res> {
  __$TherapistCopyWithImpl(this._self, this._then);

  final _Therapist _self;
  final $Res Function(_Therapist) _then;

/// Create a copy of Therapist
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? email = null,Object? availableSlotsCount = null,Object? reviewCount = null,}) {
  return _then(_Therapist(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,email: null == email ? _self.email : email // ignore: cast_nullable_to_non_nullable
as String,availableSlotsCount: null == availableSlotsCount ? _self.availableSlotsCount : availableSlotsCount // ignore: cast_nullable_to_non_nullable
as int,reviewCount: null == reviewCount ? _self.reviewCount : reviewCount // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}

// dart format on
