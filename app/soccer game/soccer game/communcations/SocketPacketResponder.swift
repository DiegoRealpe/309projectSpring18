//
//  SocketPacketResponder.swift
//  soccer game
//
//  Created by rtoepfer on 2/16/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

//contains information to map packet codes to handlers
struct PacketType {
    var dataSize : Int
    var handlerFunction : ([UInt8]) -> Void
    
    func run(data : [UInt8]){
        handlerFunction(data)
    }
}

//places recieved bytes in a queue and once enough have been recieved, executes a function with the required bytes
//packet types can be mapped through the packetType Dict
class SocketPacketResponder{
    
    private let byteQueue = PacketByteQueue()
    
    private var packetCode : UInt8 = 0
    private var bytesForPacket = 0xFFFFFFFF
    
    var packetTypeDict : [UInt8:PacketType] = [:] //modifying this when part of a packet is queued may cause issues,
                                                //fixes will need to be more complex than just stashing the current PacketType
    
    func respond(data : [UInt8]){
        
        setPacketTypeOptions(packetCode: data[0])
        byteQueue.enque(data: data)
        
        //loop, and not if statement to allow multiple packets to be executed
        while byteQueue.size >= bytesForPacket {
            guard let packetType = packetTypeDict[packetCode] else{ //todo improve handling for missing data
                print("packet number missing from dictionary")
                break
            }
            
            let data = byteQueue.deque(ammount: bytesForPacket)
            packetType.run(data: data)
            
            //reset controls if more packet data is avaible
            if byteQueue.size > 0{
                setPacketTypeOptions(packetCode: data[0])
            }else{
                resetPacketOptions()
            }
        }
    }
    
    private func setPacketTypeOptions(packetCode : UInt8){
        if self.packetCode == 0 , let packetType = packetTypeDict[packetCode]{
            bytesForPacket = packetType.dataSize
            self.packetCode = packetCode
        }
    }
    
    private func resetPacketOptions(){
        packetCode = 0
        bytesForPacket = 0xFFFFFFFF
    }
}
