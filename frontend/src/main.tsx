import React from 'react';
import { createRoot } from 'react-dom/client';
import { createHashRouter, RouterProvider } from 'react-router-dom';

import routers from '@/routers';
import './style.css';

const container = document.getElementById('root');
const root = createRoot(container!);

const router = createHashRouter(routers);

root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
