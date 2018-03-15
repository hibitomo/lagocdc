package main

import (
	"context"
	"io"
	"log"

	"github.com/reiver/go-telnet"
	"github.com/coreswitch/cmd"
	pb "github.com/coreswitch/openconfigd/proto"
	"google.golang.org/grpc"
)

const (
	defaultServer = ":2650"
	LAGOCD_MODULE = "Lagocd module"
	LAGOCD_PORT   = 10484
)

type Command struct {
	Cmd  int
	Path []string
}

func LagoDSL(command Command) error {
	if command.Cmd == cmd.Set {
		log.Println("[cmd] add", command.Path)
	} else {
		log.Println("[cmd] del", command.Path)
	}
	ret, fn, args, _ := Parser.ParseCmd(command.Path)
	if ret == cmd.ParseSuccess {
		fn.(func(int, cmd.Args) int)(command.Cmd, args)
	}
	return nil
}

// func xxx(Cmd int, Args cmd.Arg) int {
// switch Cmd {
// case cmd.Set:
// case cmd.Delete:
// }
// return cmd.Success
// }

func SetDSL() error {
	dsl, _ := SetIfDSL()
	l, _ := SetPtDSL()
	dsl = append(dsl, l...)
	l, _ = SetChDSL()
	dsl = append(dsl, l...)
	l, _ = SetConDSL()
	dsl = append(dsl, l...)
	l, _ = SetBrDSL()
	dsl = append(dsl, l...)

	telnet.DialToAndCall("localhost:12345",client{DSL: dsl})
	for _,str := range dsl {
		log.Println(str)
	}
	return nil
}

var Parser *cmd.Node

func InitAPI() {
	Parser = cmd.NewParser()
	Parser.InstallCmd([]string{"interface", "WORD", "type", "WORD"}, ConfigInterfaceType)
	Parser.InstallCmd([]string{"interface", "WORD", "device", "WORD"}, ConfigInterfaceDevice)
	Parser.InstallCmd([]string{"port", "WORD", "interface", "WORD"}, ConfigPortInterface)
	Parser.InstallCmd([]string{"channel", "WORD", "dst-addr", "WORD"}, ConfigChannelDSTAddr)
	Parser.InstallCmd([]string{"channel", "WORD", "dst-port", "WORD"}, ConfigChannelDSTPort)
	Parser.InstallCmd([]string{"channel", "WORD", "protocol", "WORD"}, ConfigChannelDSTProtocol)
	Parser.InstallCmd([]string{"controller", "WORD", "channel", "WORD"}, ConfigConChannel)
	Parser.InstallCmd([]string{"controller", "WORD", "role", "WORD"}, ConfigConRole)
	Parser.InstallCmd([]string{"controller", "WORD", "connection-type", "WORD"}, ConfigConConnectType)
	Parser.InstallCmd([]string{"bridge", "WORD", "dpid", "WORD"}, ConfigBrDpid)
	Parser.InstallCmd([]string{"bridge", "WORD", "controller", "LINE"}, ConfigBrCon)
	Parser.InstallCmd([]string{"bridge", "WORD", "port", "WORD", "port-id", "WORD"}, ConfigBrPort)
	Parser.InstallCmd([]string{"bridge", "WORD", "fail-mode", "WORD"}, ConfigBrFailMode)
}

func grpcSubscribe(conn *grpc.ClientConn, request *pb.ConfigRequest) (pb.Config_DoConfigClient, error) {
	client := pb.NewConfigClient(conn)
	stream, err := client.DoConfig(context.Background())
	if err != nil {
		return nil, err
	}
	if err = stream.Send(request); err != nil {
		return nil, err
	}
	return stream, nil
}

func main() {
	InitAPI()
	log.Println("Starting lagolog module")

	path := []string{"interface", "port", "channel", "controller", "bridge"}
	test_request := &pb.ConfigRequest{
		Type:   pb.ConfigType_SUBSCRIBE_MULTI,
		Module: LAGOCD_MODULE,
		Port:   uint32(LAGOCD_PORT),
		Path:   path,
	}

	conn, err := grpc.Dial(
		defaultServer,
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}

	// Show

	// Config
	stream, err := grpcSubscribe(conn, test_request)
	if err != nil {
		return
	}
	validating := false
	for {
		conf, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		switch conf.Type {
		case pb.ConfigType_VALIDATE_START:
			log.Println("VALIDATE_START:", conf.Path)
			validating = true
		case pb.ConfigType_VALIDATE_END:
			log.Println("VALIDATE_END:", conf.Path)
			test_request.Type = pb.ConfigType_VALIDATE_SUCCESS
			err = stream.Send(test_request)
			if err != nil {
				log.Println(err)
			}
			validating = false
		case pb.ConfigType_COMMIT_START:
			log.Println("COMMIT_START:", conf.Path)
		case pb.ConfigType_COMMIT_END:
			log.Println("COMMIT_END:", conf.Path)
			SetDSL()
		case pb.ConfigType_SET, pb.ConfigType_DELETE:
			command := Command{int(conf.Type), conf.Path}
			if validating {
			} else {
				err := LagoDSL(command)
				if err != nil {
					log.Println(err)
				}
			}
		default:
		}
	}
}
