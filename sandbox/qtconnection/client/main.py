#!/usr/bin/python

import sys, json

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
    def __init__(self, fn: str, ln: str) -> None:
        self.firstName = fn
        self.lastName = ln
        self.myarray = [10, 2, 32]
        self.mymap = { 0: False, 1: True }
        

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
        n = self.conn.bytesAvailable()
        buf = str(self.conn.read(n))
        self.readLabel.setText(buf)
        
    def _connectPressed(self):
        self.conn.connectToHost("127.0.0.1", 1200)
        
    def _disconnectPressed(self):
        self.conn.close()
        
    def _actionPressed(self):
        person = Person("Bob", "Winner")
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
