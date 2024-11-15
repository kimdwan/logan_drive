import { createContext } from "react";
import { BrowserRouter as Routers, Routes, Route } from "react-router-dom";
import { Content, Main, SignUp } from "./pkgs";
import { LoadComputerNumber } from "./settings";

export const MainContext = createContext()

function App() {
  const { computerNumber, setComputerNumber } = LoadComputerNumber()

  return (
    <div className="App">
      <MainContext.Provider value={{ computerNumber, setComputerNumber }}>
        <Routers>

          <Routes>
            <Route path = "/" element = {<Main />} />
            <Route path = "/signup/*" element = {<SignUp />} />
            <Route path = "/main/*" element = {<Content />} />
          </Routes>

        </Routers>
      </MainContext.Provider>
    </div>
  );
}

export default App;
