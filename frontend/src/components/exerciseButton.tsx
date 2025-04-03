// TODO: Create the modals
import "../App.css";
import { useState, Dispatch } from 'react';
import { Data } from '../components/appLayout'
import Modal from ".//modal";

function ExerciseButton(
  props: { create: boolean, name?: string, id?: string, fetchApi?: (p: Dispatch<Data | undefined>) => {}, setData?: Dispatch<Data | undefined> }
) {
  const [modal, setModal] = useState<boolean>(false);
  async function modalAction(formData: any) {
    if (!props.create) {

      const weight = Number(formData.get("weight"));
      const sets = Number(formData.get("sets"));
      const reps = Number(formData.get("reps"));

      if (isNaN(weight) || isNaN(sets) || isNaN(reps)) {
        // TODO: Set up some type of error thing here that says must be a number
        setModal(false);
        return;
      }

      const d = {
        create: false,
        id: props.id,
        weight: weight,
        sets: sets,
        reps: reps,
      };

      await fetch("http://localhost:8080/app", {
        method: "POST",
        body: JSON.stringify(d),
        credentials: "include"
      });

      setModal(false);
      props.fetchApi!(props.setData!);
      return;
    }

    const name = String(formData.get("name"));
    if (!(typeof name === "string")) {
      // TODO: Set up some type of error thing here that says must be a number
      setModal(false);
      return;
    }

    const d = {
      create: true,
      name: name,
    };
    await fetch("http://localhost:8080/app", {
      method: "POST",
      body: JSON.stringify(d),
      credentials: "include"
    });

    setModal(false);
    props.fetchApi!(props.setData!);
  }

  return (
    <>
      <span onClick={() => setModal(true)} className="exerciseButton">
        <h1>
          {props.create ? "CREATE NEW" : props.name?.toUpperCase()}
        </h1>
      </span>

      {modal ?
        <Modal formAction={modalAction} id={props.id} name={props.name}>
          {props.create ? (
            <>
              <h1>CREATE EXERCISE</h1>
              <label className="modal-item-label" htmlFor="name">EXERCISE NAME</label>
              <input className="modal-item-input" type="text" id="name" name="name" required />
              <div className="modal-item">
                <input className="modal-item-button" type="submit" value="SUBMIT" />
                <button className="modal-item-button" type="button" onClick={() => setModal(false)}>CLOSE</button>
              </div>
            </>
          ) : (
            <>
              <h1>{props.name?.toUpperCase()}</h1>
              <label className="modal-item-label" htmlFor="weight">WEIGHT</label>
              <input className="modal-item-input" type="text" id="weight" name="weight" required />
              <label className="modal-item-label" htmlFor="sets">SETS</label>
              <input className="modal-item-input" type="text" id="sets" name="sets" required />
              <label className="modal-item-label" htmlFor="reps">REPS</label>
              <input className="modal-item-input" type="text" id="reps" name="reps" required />
              <div className="modal-item">
                <input className="modal-item-button" type="submit" value="SUBMIT" />
                <button className="modal-item-button" type="button" onClick={() => setModal(false)}>CLOSE</button>
              </div>
            </>
          )}

        </Modal >
        : <></>
      }
    </>
  );
}

export default ExerciseButton;
