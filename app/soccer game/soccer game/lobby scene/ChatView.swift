//
//  ChatView.swift
//  soccer game
//
//  Created by rtoepfer on 4/2/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit

class ChatView: UIView, UITableViewDataSource, UITextFieldDelegate {

    var messages : [ChatMessage] = []
    
    @IBOutlet var messageTable : UITableView!
    @IBOutlet var textInput : UITextField!
    @IBOutlet var sendButton : UIButton!
    
    @IBOutlet var player0Label : UILabel!
    @IBOutlet var player1Label : UILabel!
    @IBOutlet var player2Label : UILabel!
    @IBOutlet var player3Label : UILabel!
    
    var playerLabelArray : [UILabel]!
    
    @IBOutlet var player0Emoji : EmojiTextField!
    @IBOutlet var player1Emoji : EmojiTextField!
    @IBOutlet var player2Emoji : EmojiTextField!
    @IBOutlet var player3Emoji : EmojiTextField!
    
    var playerEmojiArray : [EmojiTextField]!
    
    var onNewMessage: ( (String) -> () )?
    var onEmojiChange: ( (Int,String) -> () )?
    
    var size = 0
    static let defaultEmoji = "🐱"
    
    override init(frame: CGRect) {
        super.init(frame: frame)
    }
    
    required init?(coder aDecoder: NSCoder) {
        super.init(coder: aDecoder)
    }
    
    func labelForPlayer(_ num : Int) -> UILabel?{
        switch num {
        case 0 :
            return self.player0Label
        case 1 :
            return self.player1Label
        case 2 :
            return self.player2Label
        case 3 :
            return self.player3Label
        default :
            return nil
        }
    }
    
    func emojiForPlayer(_ num : Int) -> EmojiTextField?{
        switch num {
        case 0 :
            return self.player0Emoji
        case 1 :
            return self.player1Emoji
        case 2 :
            return self.player2Emoji
        case 3 :
            return self.player3Emoji
        default :
            return nil
        }
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
        
        for i in 0..<4{
            self.emojiForPlayer(i)?.isHidden = true
            self.labelForPlayer(i)?.isHidden = true
        }
        self.size = 0
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
    
    func addPlayer(playerNum : Int, username : String, emojiEditable : Bool){
        self.size += 1
        
        let label = self.labelForPlayer(2 * playerNum)!//skip from 0 to 2 until we are ready for 4 players
        let emoji = self.emojiForPlayer(2 * playerNum)!
        
        label.isHidden = false
        label.text = username
        
        emoji.isHidden = false
        emoji.isUserInteractionEnabled = emojiEditable
        emoji.text = ChatView.defaultEmoji
    }
    
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        self.endEditing(true)
        return false
    }
    
    func textField(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool {
        if textField == self.textInput{
            return allowTextInput(textField, shouldChangeCharactersIn: range, replacementString: string)
        }else{
            return allowEmojiInput(textField, shouldChangeCharactersIn: range, replacementString: string)
        }
    }
    
    func allowTextInput(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool{
        if string.count + textField.text!.count <= 100 {
            
            return true
        }else if string.utf8.count + textField.text!.utf8.count <= 400{
            return false
        }
        
        return false
    }
    
    func changeEmoji(playerNumber: Int,emoji: String){
        self.emojiForPlayer(playerNumber * 2)?.text = emoji
    }
    
    func allowEmojiInput(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool{
        if string.count <= 1 {
            if string.count == 1 {
                self.endEditing(true)
            }
            tellLobbySceneAboutEmojiChange(textField: textField, string)
            textField.text = string //we need to do this so we can close the keyboard and change text
            return true
        }
        
        return false
    }
    
    func tellLobbySceneAboutEmojiChange(textField: UITextField,_ emoji : String){
        for i in 0..<4{
            if self.emojiForPlayer(i) == textField {
                self.onEmojiChange?(i/2,emoji)
                return
            }
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
