import 'dart:convert';
import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:physiolink_app/features/appointments/data/appointment_repository.dart';

/// A [HttpClientAdapter] that replays a canned response, so repository
/// behaviour can be asserted without a live backend.
class _FakeAdapter implements HttpClientAdapter {
  _FakeAdapter(this.body, {this.statusCode = 200});

  final String body;
  final int statusCode;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async {
    return ResponseBody.fromString(
      body,
      statusCode,
      headers: {
        Headers.contentTypeHeader: [Headers.jsonContentType],
      },
    );
  }

  @override
  void close({bool force = false}) {}
}

AppointmentRepository _repoReturning(String body, {int statusCode = 200}) {
  final dio = Dio(BaseOptions(baseUrl: 'http://test.invalid/api'))
    ..httpClientAdapter = _FakeAdapter(body, statusCode: statusCode);
  return AppointmentRepository(dio);
}

void main() {
  group('AppointmentRepository.bookAppointment', () {
    // The endpoint returns {"id": "<uuid>"}, not an Appointment. Parsing it as
    // an Appointment used to throw, and the throw was swallowed, so the UI
    // reported "Appointment booked!" for bookings that never happened.
    test('returns the id from the {"id": ...} response', () async {
      final repo = _repoReturning('{"id":"11111111-1111-1111-1111-111111111111"}');

      final id = await repo.bookAppointment('slot-1');

      expect(id, '11111111-1111-1111-1111-111111111111');
    });

    test('accepts an _id key as well', () async {
      final repo = _repoReturning('{"_id":"22222222-2222-2222-2222-222222222222"}');

      expect(
        await repo.bookAppointment('slot-1'),
        '22222222-2222-2222-2222-222222222222',
      );
    });

    test('throws when the body carries no id', () async {
      final repo = _repoReturning('{"ok":true}');

      await expectLater(
        repo.bookAppointment('slot-1'),
        throwsA(isA<StateError>()),
      );
    });

    test('propagates a 409 conflict instead of reporting success', () async {
      final repo = _repoReturning(
        jsonEncode({'msg': 'conflict'}),
        statusCode: 409,
      );

      await expectLater(
        repo.bookAppointment('slot-1'),
        throwsA(isA<DioException>()),
      );
    });
  });
}
