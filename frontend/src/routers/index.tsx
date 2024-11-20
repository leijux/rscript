import { RouteObject } from 'react-router-dom';

import App from '@/App';
import HomePage from '@/pages/HomePage';

const routes: RouteObject[] = [
  {
    path: '/',
    element: <App />,
    children: [
      {
        path: 'home',
        element: <HomePage />,
      }
    ],
  },
];

export default routes;
