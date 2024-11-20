import React from 'react';
import { EventsOn } from '@wailsjs/runtime';

import HomePage from '@/pages/HomePage';
import { toast } from '@/components/ui/use-toast';
import { Toaster } from '@/components/ui/toaster';

EventsOn('err_msg', (err?: string) => {
  toast({
    title: 'err',
    description: err,
  });
});

function App() {
  return (
    <div className="select-none">
      <HomePage />
      <Toaster />
    </div>
  );
}

export default App;
