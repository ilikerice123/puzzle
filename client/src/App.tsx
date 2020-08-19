import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
} from "react-router-dom";
import './App.scss';
import CreatePuzzle from './components/CreatePuzzle'
import Puzzle from './components/Puzzle'
import { UserObject } from './components/game';
import Identity from './components/Identity'

function App() {
  const [user, setUser] = React.useState<UserObject | null>(null)
  return (
    <Router>
      <div>
        <nav className="navbar">
          <ul>
            <li className="nav-left title">
              <Link className="link" to="/">PUZZLE</Link>
            </li>
            <li className="nav-left item">
              <Link className="link" to="/puzzles">Create Puzzle</Link>
            </li>
            <li className="nav-right item">
              <div className="link">
                <Identity user={user} changeUser={setUser}/>
              </div>
            </li>
          </ul>
        </nav>

        <Switch>
          <Route exact path="/puzzles">
              <CreatePuzzle />
          </Route>
          <Route path="/puzzles/:id">
              <Puzzle user={user}/>
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
