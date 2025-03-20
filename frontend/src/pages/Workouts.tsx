import { useState, useEffect, Dispatch } from 'react';
import { useData } from '../components/appLayout'

type Workout = {
  weight: number,
  sets: number,
  reps: number
  time: number
};


function Workout(props: { workout: Workout }) {
  // TODO: Format this and then app should be done 
  return (
    <>
      Weight:
      {props.workout.weight}
      Reps:
      {props.workout.reps}
      Sets:
      {props.workout.sets}
      Time:
      {props.workout.time}
    </>
  );
}
function Workouts() {
  const [data] = useData();
  const [workoutData, setWorkoutData] = useState<Workout[] | undefined>(undefined);

  const fetchWorkouts = async (id: string) => {
    console.log(id);
    const raw = await fetch("http://localhost:8080/workouts", {
      credentials: 'include',
      method: "POST",
      body: id,
    });

    const json = await raw.json();
    console.log(json);
    setWorkoutData(json.workouts);
  }
  useEffect(() => {
    fetchWorkouts(data.exercises[0].id);
  }, []);
  return (
    <>
      {data ? (
        <>
          <select onChange={(e) => fetchWorkouts(e.target.value)}>
            {data.exercises.map((w) => {
              return <option value={w.id}>{w.name}</option>
            })}

          </select>
          <div>
            {workoutData?.map((w) => {
              return <Workout workout={w} />
            })}
          </div>
        </>
      ) : (
        <div>No data loaded :(</div>)}
    </>
  );
}


export default Workouts;
