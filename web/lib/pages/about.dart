import 'package:jaspr/jaspr.dart';
import 'package:jaspr/dom.dart';

class About extends StatelessComponent {
  const About({super.key});

  @override
  Component build(BuildContext context) {
    return section([
      ol([
        li([
          h3([Component.text('Docs')]),
          Component.text('PhysioLink documentation and guides.'),
        ]),
        li([
          h3([Component.text('Support')]),
          Component.text('Need help? Contact our support team.'),
        ]),
        li([
          h3([Component.text('GitHub')]),
          a(href: 'https://github.com/divijg19/physiolink', [Component.text('View on GitHub')]),
        ]),
      ]),
    ]);
  }
}
