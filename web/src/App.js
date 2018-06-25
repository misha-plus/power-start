import React, { Component } from 'react';
import './App.css';
import 'semantic-ui-css/semantic.min.css';
import {
  Button, Card, Image, Dimmer, Loader, Icon, Header, Message, Divider, Container
} from 'semantic-ui-react'
import axios from 'axios';

import AddModal from './AddModal'
import MachineCard from './MachineCard'

const CustomMessage = ({ header, text, status, onDismiss }) => {
  let color = null
  if (status === 'success') {
    color = 'green'
  } else if (status === 'error') {
    color = 'red'
  }

  return (
    <Message
      onDismiss={onDismiss}
      {...{color}}
      icon
    >
      { status === 'pending' &&
        <Icon name='circle notched' loading />
      }
      <Message.Content>
        <Message.Header>{header}</Message.Header>
        {text}
      </Message.Content>
    </Message>
  )
}

class App extends Component {
  constructor() {
    super()
    this.state = { data: null, lastMessage: null }
    this.updateData()
    setInterval(this.updateData.bind(this), 2000)
  }

  setMessage(header, text, status) {
    this.setState({ lastMessage: { header, text, status } })
  }

  clearMessage() {
    this.setState({ lastMessage: null })
  }

  async updateData() {
    try {
      const resp = await axios.get('/api/list')
      this.setState({ data: resp.data })
    } catch (e) {
      this.setMessage(
        `Can't fetch update`,
        e.message + ' / ' + e.response.data,
        'error'
      )
    }
  }

  async start(name) {
    try {
      this.setMessage(`Sending start request for "${name}"`, null, 'pending')
      await axios.post(`/api/start/${name}`)
      this.setMessage(`Start request for "${name}" was sent`, null, 'success')
      await this.updateData()
    } catch (e) {
      this.setMessage(
        `Can't send start request for "${name}"`,
        e.message + ' / ' + e.response.data,
        'error'
      )
    }
  }

  async stop(name) {
    try {
      this.setMessage(`Sending stop request for "${name}"`, null, 'pending')
      await axios.post(`/api/stop/${name}`)
      this.setMessage(`Stop request for "${name}" was sent`, null, 'success')
      await this.updateData()
    } catch (e) {
      this.setMessage(
        `Can't send stop request for "${name}"`,
        e.message + ' / ' + e.response.data,
        'error'
      )
    }
  }

  async remove(name) {
    try {
      this.setMessage(`Removing machine "${name}"`, null, 'pending')
      await axios.post(`/api/remove/${name}`)
      this.setMessage(`Machine "${name}" removed`, null, 'success')
      await this.updateData()
    } catch (e) {
      this.setMessage(
        `Can't remove machine "${name}"`,
        e.message + ' / ' + e.response.data,
        'error'
      )
    }
  }

  async add(name, mac) {
    try {
      this.setMessage(`Saving machine "${name}"`, null, 'pending')
      await axios.post(
        `/api/add`,
        {  name, mac }
      )
      this.setMessage(`Machine "${name}" saved`, null, 'success')
      await this.updateData()
    } catch (e) {
      this.setMessage(
        `Can't save machine "${name}"`,
        e.message + ' / ' + e.response.data,
        'error'
      )
    }
  }

  render() {
    if (!this.state.data) {
      return (
        <Dimmer active>
          <Loader>
            Loading machines list...
            <p>Power-Start</p>
          </Loader>
        </Dimmer>
      )
    }

    return (
      <div>
        <div className='power-start-header'>
          <Header as='h2' icon textAlign='center'>
            <Image src='icon.png' size='massive' circular />
            <Header.Content>Power-Start</Header.Content>
          </Header>
        </div>

        {this.state.lastMessage &&
          <Container className='power-start-message'>
            <CustomMessage
              header={this.state.lastMessage.header}
              text={this.state.lastMessage.text}
              status={this.state.lastMessage.status}
              onDismiss={() => this.clearMessage()}
            />
          </Container>
        }

        <Container>
          <Button
            icon
            labelPosition='left'
            onClick={() => this.setState({ isAddModalOpen: true })}
          >
            <Icon name='add' />
            Add a machine
          </Button>
          <AddModal
            isOpen={this.state.isAddModalOpen}
            onCancel={() => this.setState({ isAddModalOpen: false })}
            onAdd={(name, mac) => {
              this.setState({ isAddModalOpen: false })
              this.add(name, mac)
            }}
          >

          </AddModal>
        </Container>

        <Divider />

        <Card.Group className='power-start-centered'>

          {this.state.data.map(machine =>
            <MachineCard
              key={machine.name}
              name={machine.name}
              mac={machine.mac}
              requests={machine.requests}
              isRunning={machine.isRunning}
              onStart={(n => () => this.start(n))(machine.name)}
              onStop={(n => () => this.stop(n))(machine.name)}
              onRemove={(n => () => this.remove(n))(machine.name)}
            />
          )}
        </Card.Group>
      </div>
    );
  }
}

export default App;
