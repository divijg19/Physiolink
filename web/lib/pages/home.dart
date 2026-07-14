import 'package:jaspr/jaspr.dart';
import 'package:jaspr/dom.dart';

class Home extends StatelessComponent {
  const Home({super.key});

  @override
  Component build(BuildContext context) {
    return div(classes: 'home-container', [
      section(classes: 'hero', [
        div(classes: 'hero-content', [
          h1([Component.text('PhysioLink')]),
          p(classes: 'subtitle', [Component.text('Connecting Patients with Physiotherapists seamlessly.')]),
          div(classes: 'cta-container', [
            a(href: '/register', classes: 'button primary', [Component.text('Get Started')]),
            a(href: '/login', classes: 'button secondary', [Component.text('Login')]),
          ]),
        ]),
      ]),

      section(classes: 'features', [
        h2([Component.text('Why Choose PhysioLink?')]),
        div(classes: 'feature-grid', [
          div(classes: 'feature-card', [
            h3([Component.text('Expert Therapists')]),
            p([Component.text('Find certified and rated physiotherapists near you.')]),
          ]),
          div(classes: 'feature-card', [
            h3([Component.text('Easy Booking')]),
            p([Component.text('Book appointments instantly with real-time availability.')]),
          ]),
          div(classes: 'feature-card', [
            h3([Component.text('Track Progress')]),
            p([Component.text('Monitor your recovery journey with digital records.')]),
          ]),
        ]),
      ]),

      section(classes: 'how-it-works', [
        h2([Component.text('How It Works')]),
        ol([
          li([
            strong([Component.text('1. Sign Up')]),
            Component.text(' - Create your free account.'),
          ]),
          li([
            strong([Component.text('2. Search')]),
            Component.text(' - Browse therapists by specialty and location.'),
          ]),
          li([
            strong([Component.text('3. Book')]),
            Component.text(' - Schedule your appointment in seconds.'),
          ]),
        ]),
      ]),
    ]);
  }
}
