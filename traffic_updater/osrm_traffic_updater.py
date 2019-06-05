import sys
import time
sys.path.append('gen-py')

from proxy import ProxyService

from thrift import Thrift
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TCompactProtocol

def main():
    # Make socket
    transport = TSocket.TSocket('localhost', 6666)

    # Buffering is critical. Raw sockets are very slow
    #transport = TTransport.TBufferedTransport(transport)
    transport = TTransport.TFramedTransport(transport)

    # Wrap in a protocol
    protocol = TCompactProtocol.TCompactProtocol(transport)

    # Create a client to use the protocol encoder
    client = ProxyService.Client(protocol)

    # Connect!
    transport.open()

    # Get flows 
    start_time = time.clock()
    #flow = client.getFlowById(991239906)
    #print flow
    flows = client.getAllFlows()
    print "flows count: " + str(len(flows))
    after_get_flow_time = time.clock() 
    print "get flows time used: " + str(after_get_flow_time - start_time) + " seconds"

    # Dump to OSRM csv format
    print "TODO: dump to osrm csv"

    end_time = time.clock()
    print "dump csv time used: " + str(end_time - after_get_flow_time) + " seconds"


    # Close!
    transport.close()


if __name__ == '__main__':
    main()
    
