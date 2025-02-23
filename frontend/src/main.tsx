import { BrowserRouter, Routes, Route } from 'react-router'
import { createRoot } from 'react-dom/client'
import AppLayout from './components/appLayout.tsx'
import App from './pages/App.tsx'
import Workouts from './pages/Workouts.tsx'
import Index from './pages/Index.tsx'

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <Routes>
      <Route index element={<Index />} />
      <Route element={<AppLayout />}>
        <Route path="app" element={<App />} />
        <Route path="workouts" element={<Workouts />} />
      </Route>
    </Routes>
  </BrowserRouter>
)
