import { useState, useEffect } from 'react';
import { useData } from '../components/appLayout'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'

import workoutsStyles from './workouts.module.css'

type Workout = {
  weight: number,
  sets: number,
  reps: number,
  time: number,
  timeString: string
};

const columnHelper = createColumnHelper<Workout>();

const columns = [
  columnHelper.accessor('weight', {
    header: () => <span>WEIGHT</span>,
    cell: info => info.renderValue()
  }),
  columnHelper.accessor('sets', {
    header: () => <span>SETS</span>,
    cell: info => info.renderValue()
  }),
  columnHelper.accessor('reps', {
    header: () => <span>REPS</span>,
    cell: info => info.renderValue()
  }),
  columnHelper.accessor('timeString', {
    header: () => <span>DATE</span>,
    cell: info => info.renderValue()?.toUpperCase()
  })
]

const fallbackData: Workout[] = [];
function Workouts() {
  const [data] = useData();
  const [workoutData, setWorkoutData] = useState<Workout[] | undefined>(undefined);
  const fetchWorkouts = async (id: string) => {
    const raw = await fetch("http://localhost:8080/workouts", {
      credentials: 'include',
      method: "POST",
      body: id,
    });

    const json = await raw.json();
    let workouts: Workout[] = json.workouts;
    // Sort so that the newer records get put at the top
    if (workouts == null) {
      setWorkoutData(undefined);
      return;
    }
    workouts.sort((a, b) => b.time - a.time)
    setWorkoutData(workouts);
  }
  const table = useReactTable({
    data: workoutData ?? fallbackData,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  useEffect(() => {
    fetchWorkouts(data.exercises[0].id);
  }, []);


  return (
    <>
      {data ? (
        <>
          <select className={workoutsStyles['select']} onChange={(e) => fetchWorkouts(e.target.value)}>
            {data.exercises.map((w) => {
              return <option key={w.id} value={w.id}>{w.name.toUpperCase()}</option>
            })}

          </select>
          {workoutData ? (<>
            {table ? (
              <table className={workoutsStyles['table']}>
                <thead>
                  {table.getHeaderGroups().map(headerGroup => (
                    <tr key={headerGroup.id}>
                      {headerGroup.headers.map(header => (
                        <th key={header.id}>
                          {header.isPlaceholder
                            ? null
                            : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                        </th>
                      ))}
                    </tr>
                  ))}
                </thead>
                <tbody>
                  {table.getRowModel().rows.map(row => (
                    <tr key={row.id}>
                      {row.getVisibleCells().map(cell => (
                        <td key={cell.id}>
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : <>
            </>}
          </>) : <>
            <div style={{display: 'flex', alignItems: 'center', justifyContent:'center'}}>
              <h1>THERE ARE NO WORKOUTS FOR THIS EXERCISE</h1>

            </div>
          </>}


        </>
      ) : (
        <div>No data loaded :(</div>)}
    </>
  );
}


export default Workouts;
