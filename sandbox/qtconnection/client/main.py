#!/usr/bin/python

import sys, json, pickle

# Check for dependacies
try:
    import PyQt6.QtCore
    import PyQt6.QtWidgets
    import PyQt6.QtGui
except Exception as e:
    print(str(e))
    sys.exit(1)

from PyQt6.QtWidgets import (
    QApplication,
    QMainWindow,
    QWidget,
    QVBoxLayout,
    QPushButton,
    QLabel,
)

from PyQt6.QtNetwork import QTcpSocket

class Person(object):
    def __init__(self,
                 FirstName: str,
                 LastName: str,
                 MyArray: list[int],
                 MyMap: dict[str, bool]) -> None:
        self.firstName = FirstName
        self.lastName = LastName
        self.myArray = MyArray
        self.myMap = MyMap

class SandBoxWidget(QWidget):
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        self.conn = QTcpSocket(self)
        self.conn.connected.connect(self._tcpConnected)
        self.conn.readyRead.connect(self._tcpRead)
        
        vlayout = QVBoxLayout(self)
        connectButton = QPushButton("Connect", self)
        disconnectButton = QPushButton("Disconnect", self)
        actionButton = QPushButton("Say hi...", self)
        self.readLabel = QLabel(self)
        
        vlayout.addWidget(connectButton)
        vlayout.addWidget(disconnectButton)
        vlayout.addWidget(actionButton)
        vlayout.addWidget(self.readLabel)
        
        connectButton.pressed.connect(self._connectPressed)
        disconnectButton.pressed.connect(self._disconnectPressed)
        actionButton.pressed.connect(self._actionPressed)
        
    def _tcpConnected(self):
        print("Connected:" + self.conn.peerAddress().toString())
        print("Port:" + str(self.conn.peerPort()))
        
    def _tcpRead(self):
        # Create person object from server json
        data = self.conn.readAll().data().decode()
        j = json.loads(data)
        person = Person(**j)
        
        print(person.myArray, person.myMap)
        
        self.readLabel.setText(data)
        
    def _connectPressed(self):
        self.conn.connectToHost("127.0.0.1", 1200)
        
    def _disconnectPressed(self):
        self.conn.close()
        
    def _actionPressed(self):
        # Send json object to server
        person = Person("Bob", "Winner", [12, 32, 54], {"1": True, "0": False})
        jsonobj = json.dumps(vars(person))
        
        buf = bytearray()
        buf.extend(map(ord, jsonobj))
        
        self.conn.write(buf)

class MainWindow(QMainWindow):
    def __init__(self, parent: QWidget | None) -> None:
        super().__init__(parent)
        
        sandbox = SandBoxWidget(self)
        self.setCentralWidget(sandbox)

if __name__ == "__main__":
    app = QApplication(sys.argv[1:])
    app.setApplicationName("Qt Client")
    
    window = MainWindow(None)
    window.show()
    
    app.exec()
