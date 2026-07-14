import 'package:jaspr/jaspr.dart';
import 'package:jaspr/dom.dart';
import 'package:jaspr_router/jaspr_router.dart';

class Header extends StatelessComponent {
  const Header({super.key});

  @override
  Component build(BuildContext context) {
    final activePath = RouteState.of(context).location;

    return header([
      nav([
        for (var route in [
          (label: 'Home', path: '/'),
          (label: 'About', path: '/about'),
        ])
          div(classes: activePath == route.path ? 'active' : null, [
            Link(to: route.path, child: Component.text(route.label)),
          ]),
        a(href: '/login', [Component.text('Login')]),
      ]),
    ]);
  }
}
