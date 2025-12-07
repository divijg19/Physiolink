import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../data/therapist_repository.dart';
import '../domain/therapist.dart';

part 'therapist_controller.g.dart';

@riverpod
Future<List<Therapist>> therapists(Ref ref) async {
  final repo = ref.watch(therapistRepositoryProvider);
  return repo.getTherapists();
}
