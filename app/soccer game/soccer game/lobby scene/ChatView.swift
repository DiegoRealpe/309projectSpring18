//
//  ChatView.swift
//  soccer game
//
//  Created by rtoepfer on 4/2/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit

class ChatView: UIView, UITableViewDataSource {
    
    var messages : [ChatMessage] = []
    
    var messageTable : UITableView!
    var textInput : UITextField!
    var sendButton : UIButton!
    
    
    override init(frame: CGRect) {
        super.init(frame: frame)
        initCommon()
    }
    
    required init?(coder aDecoder: NSCoder) {
        super.init(coder: aDecoder)
        initCommon()
    }
    
    
    func initCommon(){
        self.messageTable = self.viewWithTag(1) as! UITableView
        self.sendButton = self.viewWithTag(2) as! UIButton
        self.textInput = self.viewWithTag(3) as! UITextField
    }

    
    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let message = messages[indexPath.row]
        if message.local {
            return dequeAndSetLocalCell(tableView, message)
        }else{
            return dequeAndSetRemoteCell(tableView, message)
        }
    }
    
    private func dequeAndSetLocalCell(_ tableView: UITableView, _ message: ChatMessage) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "localMessage") as! LocalChatMessageViewCell
        cell.messageLabel.text = message.text
        return cell
    }
    
    private func dequeAndSetRemoteCell(_ tableView: UITableView, _ message: ChatMessage) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "remoteMessage") as! RemoteChatMessageViewCell
        cell.usernameLabel.text = message.username + " says:"
        cell.messageLabel.text = message.text
        return cell
    }
    
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return messages.count
    }
    
    
    @IBAction func buttonWasTouched() {
        let message = ChatMessage(text:textInput.text!,username : "Ryan",local : true)
        
        messages.append(message)
        messageTable.reloadData()
        setTableToBottom(animated: true)
        
        textInput.text = ""
    }
    
    
    func loadChat(){
        for _ in 0..<2 {
            messages.append(ChatMessage(text : "hello",username: "Ryan", local : true))
            messages.append(ChatMessage(text : "world",username: "Nolan",local : false))
        }
        
        messageTable.dataSource = self
        
        setTableToBottom(animated : false)//must be called after data is loaded
    }
    
    fileprivate func setTableToBottom(animated : Bool) {
        self.messageTable.reloadData() // To populate your tableview first
        
        let totalRowSize = self.messageTable.contentSize.height
        let tableSize = self.messageTable.frame.size.height
        
        if totalRowSize > tableSize {
            
            let indexPath = IndexPath(row: self.messages.count - 1, section: 0)
            
            self.messageTable.scrollToRow(at: indexPath, at: .bottom, animated: true)
            
        }
    }
}

struct ChatMessage{
    var text : String
    var username : String
    var local : Bool
    
    init(text : String, username : String, local : Bool){
        self.text = text
        self.username = username
        self.local = local
    }
}
