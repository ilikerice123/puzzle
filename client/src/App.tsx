import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
} from "react-router-dom";
import './App.css';
import CreatePuzzle from './components/CreatePuzzle'
import Puzzle from './components/Puzzle'

function App() {
  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/puzzles">Create Puzzle</Link>
            </li>
          </ul>
        </nav>

        <Switch>
          <Route exact path="/puzzles">
              <CreatePuzzle />
          </Route>
          <Route path="/puzzles/:id">
              <Puzzle />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
