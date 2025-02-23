import { Outlet, useOutletContext, Link } from 'react-router'
import { useState, useEffect, Dispatch } from 'react'

/*
type Workout = {
  weight: number,
  sets: number,
  reps: number
};
*/

type Exercise = {
  id: string,
  name: string,
  //workouts: Workout[]
};

export type Data = {
  nameNotValid: boolean,
  name: string,
  exercises: Exercise[]
};


export async function fetchApi(setData: Dispatch<Data | undefined>) {
  // TODO: Figure out environment variables
  const raw = await fetch("http://localhost:8080/app", {
    credentials: "include",
  });


  const json = await raw.json();
  setData(json);
}
function AppLayout() {
  const [data, setData] = useState<Data | undefined>(undefined);
  useEffect(() => {
    fetchApi(setData);
  }, []);

  return (
    <>
      <div className="navBar">
        <h1>
          {data?.name.toUpperCase() + "'S GYM JOURNAL"}
        </h1>
        <a href="http://localhost:8080/logout/google">EXIT</a>
      </div>
      <Outlet context={[data, setData]} />
      <div className="footer">
        <Link to="/app">Tracker</Link>
        <Link to="/workouts">Workouts</Link>
      </div>
    </>
  );
}

export function useData() {
  return useOutletContext<[Data, Dispatch<Data | undefined>]>();
}

export default AppLayout;
