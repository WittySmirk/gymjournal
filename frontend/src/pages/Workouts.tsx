import { useData } from '../components/appLayout'
function Workouts() {
  const [data] = useData();
  console.log(data);
  return (
    <>
      {data.name}
      Workouts page
    </>
  );
}


export default Workouts;
