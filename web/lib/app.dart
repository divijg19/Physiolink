import 'package:jaspr/jaspr.dart';
import 'package:jaspr/dom.dart';
import 'package:jaspr_router/jaspr_router.dart';

import 'components/header.dart';
import 'pages/about.dart';
import 'pages/home.dart';

class App extends StatelessComponent {
  const App({super.key});

  @override
  Component build(BuildContext context) {
    return div(classes: 'main', [
      Router(routes: [
        ShellRoute(
          builder: (context, state, child) => Component.fragment([
            const Header(),
            child,
          ]),
          routes: [
            Route(path: '/', title: 'Home', builder: (context, state) => const Home()),
            Route(path: '/about', title: 'About', builder: (context, state) => const About()),
          ],
        ),
      ]),
    ]);
  }
}
