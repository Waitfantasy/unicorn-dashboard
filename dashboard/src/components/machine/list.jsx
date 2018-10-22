import React, {Component} from 'react';
import axios from "axios";
import {Button, Card, Col, Form, Icon, Input, Row, Table, Modal} from 'antd';

const FormItem = Form.Item;

const EditableContext = React.createContext();

const EditableRow = ({form, index, ...props}) => (
    <EditableContext.Provider value={form}>
        <tr {...props} />
    </EditableContext.Provider>
);

const EditableFormRow = Form.create()(EditableRow);

class EditableCell extends React.Component {
    state = {
        editing: false,
    }

    componentDidMount() {
        if (this.props.editable) {
            document.addEventListener('click', this.handleClickOutside, true);
        }
    }

    componentWillUnmount() {
        if (this.props.editable) {
            document.removeEventListener('click', this.handleClickOutside, true);
        }
    }

    toggleEdit = () => {
        const editing = !this.state.editing;
        this.setState({editing}, () => {
            if (editing) {
                this.input.focus();
            }
        });
    }

    handleClickOutside = (e) => {
        const {editing} = this.state;
        if (editing && this.cell !== e.target && !this.cell.contains(e.target)) {
            this.save();
        }
    }

    save = () => {
        const {record, handleSave} = this.props;
        this.form.validateFields((error, values) => {
            if (error) {
                return;
            }
            this.toggleEdit();
            handleSave({...record, ...values});
        });
    }

    render() {
        const {editing} = this.state;
        const {
            editable,
            dataIndex,
            title,
            record,
            index,
            handleSave,
            ...restProps
        } = this.props;
        return (
            <td ref={node => (this.cell = node)} {...restProps}>
                {editable ? (
                    <EditableContext.Consumer>
                        {(form) => {
                            this.form = form;
                            return (
                                editing ? (
                                    <FormItem style={{margin: 0}}>
                                        {form.getFieldDecorator(dataIndex, {
                                            rules: [{
                                                required: true,
                                                message: `${title} is required.`,
                                            }],
                                            initialValue: record[dataIndex],
                                        })(
                                            <Input
                                                ref={node => (this.input = node)}
                                                onPressEnter={this.save}
                                            />
                                        )}
                                    </FormItem>
                                ) : (
                                    <div
                                        className="editable-cell-value-wrap"
                                        style={{paddingRight: 24}}
                                        onClick={this.toggleEdit}
                                    >
                                        {restProps.children}
                                    </div>
                                )
                            );
                        }}
                    </EditableContext.Consumer>
                ) : restProps.children}
            </td>
        );
    }
}

const AddForm = Form.create()(
    (props) => {
        const {visible, ip, okDisable, onCancel, onOk, onChange, form} = props;
        const {getFieldDecorator} = form;
        return (
            <Modal
                okButtonProps={{disabled: okDisable}}
                visible={visible}
                title="Add Machine to Etcd"
                onCancel={onCancel}
                onOk={onOk}
            >
                <Form layout="vertical">
                    <FormItem
                        hasFeedback={ip.hasFeedback}
                        help={ip.help}
                        validateStatus={ip.validateStatus}
                        label="ip">
                        {getFieldDecorator('title', {
                            rules: [{required: true, message: '请输入收藏的标题!'}],
                        })(
                            <Input value={ip.val} onChange={onChange}/>
                        )}
                    </FormItem>
                </Form>
            </Modal>
        );
    }
);

class MachineList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            ip: {
                val: "",
                label: "",
                help: "",
                hasFeedback: true,
                validateStatus: "",
            },
            disable: true,
            editing: false,
            visible: false,
            dataSource: []
        };
        this.handleDelete = this.handleDelete.bind(this);
        this.handleAdd = this.handleAdd.bind(this);

        this.columns = [{
            title: 'Etcd Key',
            dataIndex: 'key',
            key: 'key',
        }, {
            title: 'Machine Id',
            dataIndex: 'id',
            key: 'id',
            sorter: (a, b) => a.id - b.id,
        }, {
            title: 'Machine Ip',
            dataIndex: 'ip',
            key: 'ip',
            editable: true,
            sorter: (a, b) => {
                if (a.ip.length === b.ip.length) {
                    return a.ip.localeCompare(b.ip)
                }

                return a.ip.length - b.ip.length
            }
        }, {
            title: 'Last Timestamp',
            dataIndex: 'last_timestamp',
            key: 'last_timestamp',
            sorter: (a, b) => a.last_timestamp - b.last_timestamp,
        }, {
            title: 'Action',
            key: 'action',
            render: (text, record) => (
                <span>
                    <Button type={"danger"} onClick={this.handleDelete.bind(this, record.ip)}>Delete</Button>
                    <Button className="ant-dropdown-link">
                        More actions <Icon type="down"/>
                    </Button>
                </span>
            ),
        }];
    }


    handleDelete(ip) {
        axios.post('/api/v1/machine/delete', {
            'ip': ip
        }).then(response => {

        });
        // console.log(key)
    }

    // TODO use binary search
    deleteDataSource(ip) {
        for (let i = 0; i < this.state.dataSource.length; i++) {
            let dataSource = this.state.dataSource;
            if (this.state.dataSource[i].ip === ip) {
                dataSource.splice(i, 1);
                this.setState({
                    dataSource: dataSource
                });
            }
        }
    }

    componentDidMount() {
        axios.get("/api/v1/machine/list").then(response => {
            this.setState({
                dataSource: response.data.data.machines
            });
        }).catch(error => {
            console.log(error)
        })
    }


    handleAdd() {
        this.setState({
            visible: true,
        });
    }

    handleSave = (row) => {
        const newData = [...this.state.dataSource];
        const index = newData.findIndex(item => row.key === item.key);
        const item = newData[index];
        newData.splice(index, 1, {
            ...item,
            ...row,
        });
        this.setState({dataSource: newData});
    };


    render() {
        const components = {
            body: {
                row: EditableFormRow,
                cell: EditableCell
            },
        };

        const columns = this.columns.map((col) => {
            if (!col.editable) {
                return col;
            }
            return {
                ...col,
                onCell: record => ({
                    record,
                    editable: col.editable,
                    dataIndex: col.dataIndex,
                    title: col.title,
                    handleSave: this.handleSave
                }),
            };
        });

        return (
            <div className="gutter-example">
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="Machine List">
                                <div style={{marginBottom: "13px"}}>
                                    <Button type="primary" onClick={this.handleAdd}>Add</Button>
                                    <AddForm
                                        okDisable={this.state.disable}
                                        ip={this.state.ip}
                                        visible={this.state.visible}
                                        onCancel={() => {
                                            this.setState({
                                                visible: false,
                                            });
                                        }}

                                        onOk={() => {
                                            if (!this.state.ip.validateStatus) {
                                                return;
                                            }
                                            console.log(this.state.ip)
                                        }}

                                        onChange={(e) => {
                                            const ip = e.target.value;
                                            const reg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
                                            if (!reg.test(ip)) {
                                                this.setState({
                                                    ip: {
                                                        label: "error",
                                                        help: "input ip invalid!",
                                                        validateStatus: "error",
                                                        val: e.target.value,
                                                    },
                                                    disable: true,
                                                });
                                            } else {
                                                this.setState({
                                                    ip: {
                                                        label: "success",
                                                        validateStatus: "success",
                                                        val: e.target.value,
                                                    },
                                                    disable: false,
                                                });
                                            }
                                        }}
                                    />
                                </div>
                                <Table
                                    bordered
                                    components={components}
                                    dataSource={this.state.dataSource}
                                    columns={columns}
                                    onChange={this.handleChange}
                                />
                            </Card>
                        </div>
                    </Col>
                </Row>
            </div>
        );
    }
}

export default MachineList;