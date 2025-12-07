import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import '../domain/therapist.dart';
import '../../appointments/presentation/appointment_controller.dart';

class TherapistDetailScreen extends ConsumerWidget {
  final Therapist therapist;

  const TherapistDetailScreen({super.key, required this.therapist});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final availabilityAsync = ref.watch(
      therapistAvailabilityProvider(therapist.id),
    );

    return Scaffold(
      appBar: AppBar(title: Text(therapist.email)),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Available Slots',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 10),
            Expanded(
              child: availabilityAsync.when(
                data: (slots) {
                  if (slots.isEmpty) {
                    return const Center(child: Text('No available slots.'));
                  }
                  return ListView.builder(
                    itemCount: slots.length,
                    itemBuilder: (context, index) {
                      final slot = slots[index];
                      return Card(
                        child: ListTile(
                          title: Text(
                            DateFormat(
                              'MMM d, y - h:mm a',
                            ).format(slot.startTime.toLocal()),
                          ),
                          trailing: ElevatedButton(
                            onPressed: () async {
                              try {
                                await ref
                                    .read(
                                      appointmentControllerProvider.notifier,
                                    )
                                    .bookAppointment(slot.id);
                                if (context.mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    const SnackBar(
                                      content: Text('Appointment booked!'),
                                    ),
                                  );
                                  // Refresh availability
                                  ref.invalidate(
                                    therapistAvailabilityProvider(therapist.id),
                                  );
                                }
                              } catch (e) {
                                if (context.mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text('Failed to book: $e'),
                                    ),
                                  );
                                }
                              }
                            },
                            child: const Text('Book'),
                          ),
                        ),
                      );
                    },
                  );
                },
                loading: () => const Center(child: CircularProgressIndicator()),
                error: (err, stack) => Center(child: Text('Error: $err')),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
