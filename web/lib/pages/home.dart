import 'package:jaspr/jaspr.dart';

class Home extends StatelessComponent {
  const Home({super.key});

  @override
  Component build(BuildContext context) {
    return div(classes: 'home-container', [
      section(classes: 'hero', [
        div(classes: 'hero-content', [
          h1([text('PhysioLink')]),
          p(classes: 'subtitle', [text('Connecting Patients with Physiotherapists seamlessly.')]),
          div(classes: 'cta-container', [
            a(href: 'http://localhost:8080/register', classes: 'button primary', [text('Get Started')]),
            a(href: 'http://localhost:8080/login', classes: 'button secondary', [text('Login')]),
          ]),
        ]),
      ]),

      section(classes: 'features', [
        h2([text('Why Choose PhysioLink?')]),
        div(classes: 'feature-grid', [
          div(classes: 'feature-card', [
            h3([text('Expert Therapists')]),
            p([text('Find certified and rated physiotherapists near you.')]),
          ]),
          div(classes: 'feature-card', [
            h3([text('Easy Booking')]),
            p([text('Book appointments instantly with real-time availability.')]),
          ]),
          div(classes: 'feature-card', [
            h3([text('Track Progress')]),
            p([text('Monitor your recovery journey with digital records.')]),
          ]),
        ]),
      ]),

      section(classes: 'how-it-works', [
        h2([text('How It Works')]),
        ol([
          li([
            strong([text('1. Sign Up')]),
            text(' - Create your free account.'),
          ]),
          li([
            strong([text('2. Search')]),
            text(' - Browse therapists by specialty and location.'),
          ]),
          li([
            strong([text('3. Book')]),
            text(' - Schedule your appointment in seconds.'),
          ]),
        ]),
      ]),
    ]);
  }
}
