import { useState, useEffect } from 'react'
import { useData, fetchApi } from '../components/appLayout'
import ExerciseButton from "../components/exerciseButton";
import Modal from "../components/modal";
import '../App.css'

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

    await fetch("http://localhost:8080/app", {
      method: "POST",
      body: JSON.stringify(d),
      credentials: "include"
    });

    setWrongName("false");
    setNameNotValid(false);
    fetchApi(setData);
  }

  useEffect(() => {
    if (data != undefined && data!.nameNotValid) {
      setNameNotValid(true);
    }
  }, []);

  return (
    <>
      {
        data ? (
          <>
            {nameNotValid ? (
              <>
                <Modal formAction={nameAction} >
                  <h1>CREATE USER NAME</h1>
                  <label className="modal-item" htmlFor="name">PLEASE ENETER YOUR NAME</label>
                  <input className="modal-item" type="text" id="name" name="name" required />
                  <div className="modal-item">
                    <input type="submit" value="SUBMIT" />
                  </div>
                  <p className="modal-error">{wrongName == "string" ? "MAKE SURE NAME IS VALID TEXT" : ""}</p>
                </Modal>
              </>
            ) : (<></>)}
            <>
              <div className="exerciseButtonGrid">
                {data.exercises.map((e, k) => {
                  return <ExerciseButton key={k} create={false} name={e.name} id={e.id} />
                })}
                <ExerciseButton create={true} fetchApi={fetchApi} setData={setData} />
              </div >
            </>
          </>
        ) : (
          <div>No data loaded :(</div>
        )
      }
    </>
  )
}

export default App
