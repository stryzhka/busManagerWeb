import {useEffect} from "react";

const App_test = () => {
    useEffect(() => {
        fetch("http://localhost:8080/api/buses/")
            .then(response => response.json())
            .then(data => {
                console.log(data);
            })
            .catch(error => {
                console.error("Error:", error);
            });
    }, []);
    return (
        <div>ой извини пожалуйста я пидор</div>

    )
}

export default App_test