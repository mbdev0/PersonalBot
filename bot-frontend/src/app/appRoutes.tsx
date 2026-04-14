import { Route, Routes } from 'react-router';
import { RouteElements } from './routes';

export function AppRoutes() {
  return (
    <Routes>
      {RouteElements.map((route) => {
        const Element = route.element;
        return <Route path={route.url} element={<Element />}></Route>;
      })}
    </Routes>
  );
}
