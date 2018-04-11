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
    
    @IBOutlet var messageTable : UITableView!
    @IBOutlet var textInput : UITextField!
    @IBOutlet var sendButton : UIButton!
    
    var onNewMessage: ( (String) -> () )?
    
    
    override init(frame: CGRect) {
        super.init(frame: frame)
        initCommon()
    }
    
    required init?(coder aDecoder: NSCoder) {
        super.init(coder: aDecoder)
        initCommon()
    }
    
    
    func initCommon(){

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
        let text = textInput.text!
        guard !text.isEmpty else{
            return
        }
        
        
        let message = ChatMessage(text: text,username : "Ryan",local : true)
        
        messages.append(message)
        messageTable.reloadData()
        setTableToBottom(animated: true)
        
        textInput.text = ""
        
        self.onNewMessage?(text)
    }
    
    func addRemoteMessage(_ message : String, from : String){
        let message = ChatMessage(text: message, username : from ,local : false)
        
        messages.append(message)
        
        
        setTableToBottonInMainThread(animated: true)
    }
    
    func setTableToBottonInMainThread(animated : Bool){
        DispatchQueue.main.async {
           self.setTableToBottom(animated: animated)
        }
    }
    
    func loadChat(){
        
        messageTable.dataSource = self
        messages = []
        
        textInput.text = ""
        
        setTableToBottom(animated : false)//must be called after data is loaded
    }
    
    fileprivate func setTableToBottom(animated : Bool) {
        self.messageTable.reloadData()
        
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
