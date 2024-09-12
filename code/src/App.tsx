import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Register from "./component/Auth/Register";
import Login from "./component/Auth/Login";
import Home from "./component/Home";
import Project from "./component/Project";

function App() {
  const routes = createBrowserRouter([
    {
      path: "/",
      element: <Project />,
    },
    {
      path: "/register",
      element: <Register />,
    },
    {
      path: "/login",
      element: <Login />,
    },
  ]);
  return <RouterProvider router={routes} />;
}

export default App;
