import React, {Component} from 'react';
import {Table} from 'antd';
import axios from "axios";

class MachineList extends Component {
    constructor(props) {
        super(props);
        this.dateSource = {}
    }

    componentWillMount() {
        axios.get("/machine/list").then(response => {
            console.log(response)
        }).catch(error => {
            console.log(error)
        })
    }

    render() {
        return (
            <div>
                <Table />
            </div>
        );
    }

}

export default MachineList;