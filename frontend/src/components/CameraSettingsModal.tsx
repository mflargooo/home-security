import { useState } from "react";
import { updateCamera } from "../utils/api";
import { Modal } from "./Modal";
import type { Camera } from "../types/Camera";
import "../index.css"

type CameraSettingsModalProps = {
    camera: Camera;
    onClose: () => void;
}

export function CameraSettingsModal({ camera, onClose } : CameraSettingsModalProps) {
    const [canUpdate, setCanUpdate] = useState<boolean>(false);

    const handleFormSubmit = async (e : any) => {
        e.preventDefault();
        setCanUpdate(false);

        const formData = new FormData(e.currentTarget);
        const updatedData = Object.fromEntries(formData.entries())
        
        camera.id = updatedData["ID"].toString()
        camera.name = updatedData["Name"].toString()
        camera.rtsp_url = updatedData["RTSP Url"].toString()

        await updateCamera(camera.id, camera)
    }

    const isFormChanged = (e : any) => {
        const formData = new FormData(e.currentTarget);
        const data = Object.fromEntries(formData.entries());
        return data["Camera Name"] !== camera.name || data["RTSP Url"] !== camera.rtsp_url
    }

    return (
        <Modal className="w-full max-w-lg lg:max-w-xl" title={"SETTINGS"} onClose={() => onClose()}>
            <div className="bg-black h-[.5px]"></div>
            <form onChange={(e) => setCanUpdate(isFormChanged(e))} onSubmit={handleFormSubmit} className="p-5">
                <div className="form-category">
                    <div className="form-category-header">Camera</div>
                    <div className="form-category-divider"></div>
                    <div className="form-entry">
                        <label htmlFor="name">Name: </label>
                        <input id="name" name="Camera Name" type="text" defaultValue={camera.name} />
                    </div>
                    <div className="form-entry">
                        <label htmlFor="id">ID: </label>
                        <input id="id" name="Camera ID" type="text" readOnly defaultValue={camera.id} />
                    </div>
                    <div className="form-entry">
                        <label htmlFor="rtsl_url">RTSP Url: </label>
                        <input id="rtsp_url" name="RTSP Url" type="text" defaultValue={camera.rtsp_url} />
                    </div>
                </div>
                <div className="flex justify-end mt-[1rem] w-full h-fit">
                    <button id="submit" disabled={!canUpdate}>
                    UPDATE</button>
                </div>
            </form>
        </Modal>
    )
}