import { useState, useEffect } from 'react'
import { useData, fetchApi } from '../components/appLayout'
import ExerciseButton from "../components/exerciseButton";
import Modal from "../components/modal";
// import '../App.css'
import appStyles from './app.module.css'
import modalStyles from '../components/modal.module.css'

type wrongType = "false" | "string";

function App() {
  const [data, setData] = useData();
  const [nameNotValid, setNameNotValid] = useState<boolean>(false);
  const [wrongName, setWrongName] = useState<wrongType>("false");

  async function nameAction(formData: any) {
    const name = String(formData.get("name"));
    if (!(typeof name === "string")) {
      setWrongName("string");
      return;
    }
    const d = {
      createName: true,
      name: name,
    };

    await fetch(import.meta.env.VITE_BACKEND_URL + "/app", {
      method: "POST",
      body: JSON.stringify(d),
      credentials: "include"
    });

    setWrongName("false");
    setNameNotValid(false);
    fetchApi(setData);
  }

  useEffect(() => {
    if (data != undefined) {
      setNameNotValid(data!.nameNotValid)
    }
  }, [data]);

  return (
    <>
      {
        data ? (
          <>
            {nameNotValid ? (
              <>
                <Modal formAction={nameAction} >
                  <h1>CREATE USER NAME</h1>
                  <label className={modalStyles['item-label']} htmlFor="name">PLEASE ENTER YOUR NAME</label>
                  <input className={modalStyles['item-input']} type="text" id="name" name="name" required />
                  <div className={modalStyles['item']}>
                    <input className={modalStyles['item-button']} type="submit" value="SUBMIT" />
                  </div>
                  <p className={modalStyles['error']}>{wrongName == "string" ? "MAKE SURE NAME IS VALID TEXT" : ""}</p>
                </Modal>
              </>
            ) :
              <>
                <div className={appStyles['grid']}>
                  {data.exercises.map((e, k) => {
                    return <ExerciseButton key={k} create={false} name={e.name} id={e.id} />
                  })}
                  <ExerciseButton create={true} fetchApi={fetchApi} setData={setData} />
                </div >
              </>
            }
          </>
        ) : (
          <div>No data loaded :(</div>
        )
      }
    </>
  )
}

export default App
