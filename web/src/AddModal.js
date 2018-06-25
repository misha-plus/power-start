import React, { Component } from 'react'
import { Button, Icon, Modal, Form, Input } from 'semantic-ui-react'

class AddModal extends Component {
  constructor({ isOpen, onAdd, onCancel }) {
    super()
    this.props = { isOpen, onAdd, onCancel }
    this.state = { name: "", mac: "" }
  }

  onChange(a, b) {
    this.setState({ [b.name]: b.value })
  }

  render() {
    return (
      <div>
        <Modal
          open={this.props.isOpen}
          onClose={this.props.onCancel}
        >
          <Modal.Header>Add a machine</Modal.Header>
          <Modal.Content>
            <Modal.Description>
              <Form>
                <Form.Field>
                  <label>Name</label>
                  <Input
                    placeholder='Name'
                    name='name'
                    value={this.state.name}
                    onChange={this.onChange.bind(this)}
                  />
                </Form.Field>
                <Form.Field>
                  <label>MAC</label>
                  <Input
                    placeholder='MAC'
                    name='mac'
                    value={this.state.mac}
                    onChange={this.onChange.bind(this)}
                  />
                </Form.Field>
              </Form>
            </Modal.Description>
          </Modal.Content>
          <Modal.Actions>
            <Button
              negative
              onClick={this.props.onCancel}
            >
              Cancel
            </Button>
            <Button
              positive
              icon='checkmark'
              labelPosition='right'
              content='Add'
              onClick={() => {
                this.props.onAdd(this.state.name, this.state.mac)
              }}
            />
          </Modal.Actions>
        </Modal>
      </div>
    )
  }
}

export default AddModal
