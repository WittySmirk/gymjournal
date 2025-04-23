import { Dispatch, ReactNode } from "react";
import modalStyles from './modal.module.css'
function Modal(props: { children: ReactNode, formAction: (formData: any) => void, id?: string, name?: string, open?: Dispatch<boolean> }) {
  async function modalAction(formData: any) {
    props.formAction(formData);
  }
  return (
    <div className={modalStyles['modal']}>
      {/* @ts-ignore */}
      <form className={modalStyles['form']} action={modalAction}>
        {props.children}
      </form >
    </div >
  );
}

export default Modal;
