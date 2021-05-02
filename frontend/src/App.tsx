import {
  BrowserRouter,
  Switch,
  Route,
} from "react-router-dom";
import Home from "./pages/Home";
import Login from "./pages/Login";
import "./assets/scss/style.scss";
import "bootstrap-icons/font/bootstrap-icons.css";

const App = () : JSX.Element => {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path="/">
          <Home />
        </Route>
        {/* Authentication page */}
        <Route path="/auth/login">
          <Login />
        </Route>
      </Switch>
    </BrowserRouter>
  );
}

export default App;
